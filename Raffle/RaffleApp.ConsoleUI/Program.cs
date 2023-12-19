using RaffleApp.Core;

var service = new RaffleService();

while (true)
{
    Console.Clear();
    Console.WriteLine("Welcome to My Raffle App");
    Console.WriteLine($"Status: {GetRaffleStatus(service)}\n");
    Console.WriteLine("[1] Start a New Draw");
    Console.WriteLine("[2] Buy Tickets");
    Console.WriteLine("[3] Run Raffle");
    Console.WriteLine();

    Console.Write("Enter your choice: ");
    var choice = Console.ReadLine();

    switch (choice)
    {
        case "1":
            Console.Clear();
            service.StartNewDraw();
            Console.Write("Press any key to return to the main menu");
            Console.ReadKey();
            break;

        case "2":
            Console.Clear();
            service.BuyTickets();
            Console.Write("Press any key to return to the main menu");
            Console.ReadKey();
            break;

        case "3":
            Console.Clear();
            service.RunRaffle();
            Console.Write("Press any key to return to the main menu");
            Console.ReadKey();
            break;

        default:
            Console.WriteLine("Invalid choice. Please enter 1, 2, or 3.");
            break;
    }
}


static string GetRaffleStatus(IRaffleService service)
{
    return !service.Open
        ? "Draw has not started"
        : $"Draw is ongoing. Raffle pot size is ${service.PotSize}";
}