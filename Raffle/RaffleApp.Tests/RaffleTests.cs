using RaffleApp.Core;

namespace RaffleApp.Tests;

[TestFixture]
public class RaffleTests
{
    [Test]
    public void StartNewDraw_ShouldSetOpenToTrueAndIncreasePotSize()
    {
        // Arrange
        var raffle = new Raffle();

        // Act
        raffle.StartNewDraw();
        
        // Assert
        Assert.Multiple(() =>
        {
            Assert.That(raffle.Open);
            Assert.That(raffle.PotSize, Is.EqualTo(Constants.PotSeed));
        });
    }

    [Test]
    public void BuyTickets_ShouldAddUserAndTickets()
    {
        // Arrange
        var raffle = new Raffle();
        raffle.StartNewDraw();

        // Act
        var result = raffle.BuyTickets("John", 3);
        
        // Assert
        Assert.Multiple(() =>
        {
            Assert.That(result.Item1.Name, Is.EqualTo("John"));
            Assert.That(result.Item2, Has.Count.EqualTo(3));
            Assert.That(raffle.PotSize, Is.EqualTo(Constants.PotSeed + 15));
        });
    }

    [Test]
    public void BuyTickets_WhenDrawNotStarted_ShouldThrowException()
    {
        // Arrange
        var raffle = new Raffle();

        // Act and Assert
        Assert.Throws<Exception>(() => raffle.BuyTickets("Alice", 2), "Draw has not started");
    }

    [Test]
    public void BuyTickets_WhenExceedingMaxTicketsPerUser_ShouldThrowException()
    {
        // Arrange
        var raffle = new Raffle();
        raffle.StartNewDraw();
        raffle.BuyTickets("Bob", Constants.MaxTicketsPerUser - 2);

        // Act and Assert
        Assert.Throws<Exception>(() => raffle.BuyTickets("Bob", 3),
            "You can only purchase 2 more tickets in this draw.");
        Assert.That(raffle.Users["Bob"].Tickets, Has.Count.EqualTo(Constants.MaxTicketsPerUser - 2));
    }
    
    [Test]
    public void RunRaffle_WhenDrawNotStarted_ShouldThrowException()
    {
        // Arrange
        var raffle = new Raffle();

        // Act and Assert
        Assert.Throws<Exception>(() => raffle.RunRaffle(), "Draw has not started");
    }

    [Test]
    public void RunRaffle_ShouldSetWinnerAndCalculateRewards()
    {
        // Arrange
        var raffle = new Raffle(rand: new Random(3));
        raffle.StartNewDraw();
        raffle.BuyTickets("User1", 2);
        raffle.BuyTickets("User2", 2);

        // Act
        raffle.RunRaffle();
        
        // Assert
        Assert.Multiple(() =>
        {
            Assert.That(raffle.Winner, Is.Not.Null);
            Assert.That(raffle.Group2Winners, Has.Count.EqualTo(2));
            Assert.That(raffle.Group3Winners, Has.Count.EqualTo(0));
            Assert.That(raffle.Group4Winners, Is.Empty);
            Assert.That(raffle.Group5Winners, Is.Empty);
            Assert.That(raffle.PotSize, Is.EqualTo(108));
            Assert.That(raffle.Open, Is.False);
        });
    }

    [Test]
    public void RunRaffle_ShouldDistributeRewardsCorrectly()
    {
        // Arrange
        var raffle = new Raffle(rand: new Random(3));
        raffle.StartNewDraw();
        raffle.BuyTickets("User1", 2);
        raffle.BuyTickets("User2", 2);

        // Act
        raffle.RunRaffle();
        
        // Assert
        Assert.Multiple(() =>
        {
            Assert.That(raffle.PotSize, Is.EqualTo(108));
            Assert.That(raffle.Group2Winners[raffle.Users["User1"]].Item2, Is.EqualTo(6));
            Assert.That(raffle.Group2Winners[raffle.Users["User2"]].Item2, Is.EqualTo(6));
        });
    }
}