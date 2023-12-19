namespace RaffleApp.Core;

public class Raffle
{
    private readonly Random _rand = new();
    private readonly Rewards _rewards = new();

    public bool Open { get; private set; }
    public double PotSize { get; private set; }
    public Ticket? Winner { get; set; }
    public Dictionary<string, User> Users { get; set; } = new();

    public Dictionary<User, (List<Ticket>, double)> Group2Winners { get; private set; } = new();

    public Dictionary<User, (List<Ticket>, double)> Group3Winners { get; private set; } = new();

    public Dictionary<User, (List<Ticket>, double)> Group4Winners { get; private set; } = new();

    public Dictionary<User, (List<Ticket>, double)> Group5Winners { get; private set; } = new();

    public Raffle()
    {
        
    }
    
    public Raffle(Random rand)
    {
        _rand = rand;
    }

    public void StartNewDraw()
    {
        if (Open) return;
        PotSize += Constants.PotSeed;
        Open = true;
        Users = new Dictionary<string, User>();
    }

    public (User, List<Ticket>) BuyTickets(string name, int numTickets)
    {
        if (!Open)
        {
            throw new Exception("Draw has not started");
        }

        if (!Users.TryGetValue(name, out var user))
        {
            user = new User(name, []);
            Users.Add(name, user);
        }

        if (user.Tickets.Count + numTickets > Constants.MaxTicketsPerUser)
        {
            var remaining = Constants.MaxTicketsPerUser - user.Tickets.Count;
            if (remaining == 0)
            {
                throw new Exception(
                    $"You have already purchased the maximum number of tickets ({Constants.MaxTicketsPerUser}) in this draw.");
            }

            throw new Exception(
                $"You can only purchase {remaining} more ticket{(remaining > 1 ? "s" : "")} in this draw.");
        }

        var newTickets = Enumerable.Range(1, numTickets)
            .Select(_ => GenerateTicket())
            .ToList();
        user.Tickets.AddRange(newTickets);
        PotSize += numTickets * Constants.TicketPrice;
        return (user, newTickets);
    }

    public void RunRaffle()
    {
        if (!Open)
        {
            throw new Exception("Draw has not started");
        }

        Winner = GenerateTicket();
        (Group2Winners, Group3Winners, Group4Winners, Group5Winners) = GetWinners();
        double totalRewards = 0;
        totalRewards += CalculateRewards(Group2Winners, _rewards.Group2) +
                        CalculateRewards(Group3Winners, _rewards.Group3) +
                        CalculateRewards(Group4Winners, _rewards.Group4) +
                        CalculateRewards(Group5Winners, _rewards.Group5);
        PotSize -= totalRewards;
        Open = false;
    }

    private Ticket GenerateTicket()
    {
        var uniqueNumbers = new List<int>();
        var usedNumbers = new HashSet<int>();

        while (uniqueNumbers.Count < 5)
        {
            var num = _rand.Next(1, 16);

            if (usedNumbers.Add(num))
            {
                uniqueNumbers.Add(num);
            }
        }

        return new Ticket(uniqueNumbers);
    }

    private (
        Dictionary<User, (List<Ticket>, double)>,
        Dictionary<User, (List<Ticket>, double)>,
        Dictionary<User, (List<Ticket>, double)>,
        Dictionary<User, (List<Ticket>, double)>) GetWinners()
    {
        var group2Winners = new Dictionary<User, (List<Ticket>, double)>();
        var group3Winners = new Dictionary<User, (List<Ticket>, double)>();
        var group4Winners = new Dictionary<User, (List<Ticket>, double)>();
        var group5Winners = new Dictionary<User, (List<Ticket>, double)>();

        foreach (var user in Users.Values)
        {
            foreach (var ticket in user.Tickets)
            {
                var matchedCount = CountMatchingNumbers(ticket.Numbers, Winner!.Numbers);
                switch (matchedCount)
                {
                    case 2:
                        if (!group2Winners.TryAdd(user, ([ticket], 0)))
                        {
                            group2Winners[user].Item1.Add(ticket);
                        }

                        break;
                    case 3:
                        if (!group3Winners.TryAdd(user, ([ticket], 0)))
                        {
                            group3Winners[user].Item1.Add(ticket);
                        }

                        break;
                    case 4:
                        if (!group4Winners.TryAdd(user, ([ticket], 0)))
                        {
                            group4Winners[user].Item1.Add(ticket);
                        }

                        break;
                    case 5:
                        if (!group5Winners.TryAdd(user, ([ticket], 0)))
                        {
                            group5Winners[user].Item1.Add(ticket);
                        }

                        break;
                }
            }
        }

        return (group2Winners, group3Winners, group4Winners, group5Winners);
    }

    private static int CountMatchingNumbers(IEnumerable<int> ticketNumbers, List<int> winningNumbers)
    {
        var winningNumCount = new Dictionary<int, int>(winningNumbers.Count);

        foreach (var num in winningNumbers)
        {
            if (!winningNumCount.TryGetValue(num, out var value))
            {
                winningNumCount[num] = 1;
            }
            else
            {
                winningNumCount[num] = ++value;
            }
        }

        return ticketNumbers.Count(num => winningNumCount.TryGetValue(num, out var count) && count > 0);
    }

    private double CalculateRewards(IDictionary<User, (List<Ticket>, double)> winners,
        double rewardPercentage)
    {
        var totalTickets = winners.Values.Sum(tickets => tickets.Item1.Count);
        double totalRewards = 0;

        foreach (var user in winners.Keys)
        {
            var reward = rewardPercentage * PotSize / totalTickets * winners[user].Item1.Count;
            totalRewards += reward;
            winners[user] = (winners[user].Item1, reward);
        }

        return totalRewards;
    }
}