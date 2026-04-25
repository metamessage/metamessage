namespace MetaMessage.Mm;

public static class TimeUtil
{
    private static readonly DateTime Epoch = new DateTime(1970, 1, 1, 0, 0, 0, DateTimeKind.Utc);

    public static long EpochSeconds(DateTime dateTime)
    {
        if (dateTime.Kind != DateTimeKind.Utc)
        {
            dateTime = dateTime.ToUniversalTime();
        }
        return (long)(dateTime - Epoch).TotalSeconds;
    }

    public static long DaysSinceEpochUtc(DateTime dateTime)
    {
        if (dateTime.Kind != DateTimeKind.Utc)
        {
            dateTime = dateTime.ToUniversalTime();
        }
        return (long)(dateTime - Epoch).TotalDays;
    }

    public static long SecondsOfDay(DateTime dateTime)
    {
        return (long)(dateTime.Hour * 3600 + dateTime.Minute * 60 + dateTime.Second);
    }

    public static DateTime FromEpochSeconds(long seconds)
    {
        return Epoch.AddSeconds(seconds).ToLocalTime();
    }

    public static DateTime FromDaysSinceEpoch(long days)
    {
        return Epoch.AddDays(days).ToLocalTime();
    }

    public static DateTime FromSecondsOfDay(long seconds)
    {
        int hours = (int)(seconds / 3600);
        int minutes = (int)((seconds % 3600) / 60);
        int secs = (int)(seconds % 60);
        return new DateTime(1970, 1, 1, hours, minutes, secs, DateTimeKind.Local);
    }
}