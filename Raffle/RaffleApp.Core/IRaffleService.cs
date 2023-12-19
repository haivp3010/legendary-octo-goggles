namespace RaffleApp.Core;

public interface IRaffleService
{
    void StartNewDraw();
    void BuyTickets();
    void RunRaffle();
    bool Open { get; }
    double PotSize { get; }
}