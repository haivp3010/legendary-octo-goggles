namespace RaffleApp.Core;

using System;

public class RaffleService : IRaffleService
{
    private readonly Raffle _raffle = new();

    public bool Open => _raffle.Open;
    public double PotSize => _raffle.PotSize;

    public void StartNewDraw()
    {
        _raffle.StartNewDraw();
        Console.WriteLine($"New Raffle draw has been started. Initial pot size: ${_raffle.PotSize}");
    }

    public void BuyTickets()
    {
        var input = GetUserInput("Enter your name, number of tickets to purchase (e.g., James,1): ");
        var parts = input.Split(",");
        if (parts.Length != 2 || !int.TryParse(parts[1], out var numTickets))
        {
            Console.WriteLine("Invalid input format. Please enter name and number of tickets.");
            return;
        }

        var name = parts[0].Trim();
        try
        {
            var (user, tickets) = _raffle.BuyTickets(name, numTickets);
            Console.WriteLine($"Hi {user.Name}, you have purchased {numTickets} ticket{(numTickets > 1 ? "s" : "")}");
            for (var i = 0; i < tickets.Count; i++)
            {
                Console.WriteLine($"Ticket {i + 1}: {tickets[i].PrintNumbers()}");
            }
            Console.WriteLine();
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }
    }

    public void RunRaffle()
    {
        Console.WriteLine("Running Raffle..");
        try
        {
            _raffle.RunRaffle();
            Console.WriteLine($"Winning Ticket is {_raffle.Winner!.PrintNumbers()}");
            DisplayWinners("Group 2 Winners", _raffle.Group2Winners);
            DisplayWinners("Group 3 Winners", _raffle.Group3Winners);
            DisplayWinners("Group 4 Winners", _raffle.Group4Winners);
            DisplayWinners("Group 5 Winners (Jackpot)", _raffle.Group5Winners);
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }
    }

    private static void DisplayWinners(string groupName, IDictionary<User, (List<Ticket>, double)> winners)
    {
        Console.WriteLine($"{groupName}:");

        if (winners.Count == 0)
        {
            Console.WriteLine("Nil");
        }

        foreach (var user in winners.Keys)
        {
            Console.WriteLine(
                $"{user.Name} with {winners[user].Item1.Count} winning ticket(s) - ${winners[user].Item2}");
        }

        Console.WriteLine();
    }

    private static string GetUserInput(string prompt)
    {
        Console.Write(prompt);
        return Console.ReadLine()!;
    }
}