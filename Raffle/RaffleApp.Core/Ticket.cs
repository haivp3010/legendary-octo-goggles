namespace RaffleApp.Core;

public record Ticket(List<int> Numbers)
{
    public string PrintNumbers() => string.Join(" ", Numbers.Select(num => num.ToString()));
}
