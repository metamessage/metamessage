<?php

namespace io\metamessage\mm;

class TimeUtil {
    public const EPOCH = '1970-01-01 00:00:00';

    public static function zoneFromHourOffset(int $hours): string {
        if ($hours < -12 || $hours > 14) {
            throw new \Exception('Location offset hours must be between -12 and +14');
        }
        $sign = $hours >= 0 ? '+' : '-';
        $hoursAbs = abs($hours);
        return sprintf('UTC%s%02d:00', $sign, $hoursAbs);
    }

    public static function zoneOffsetHours(?string $zone): int {
        if ($zone === null) {
            return 0;
        }
        $dt = new \DateTime(self::EPOCH, new \DateTimeZone($zone));
        return $dt->getOffset() / 3600;
    }

    public static function daysSinceEpochUtc(\DateTimeInterface $date): int {
        $epoch = new \DateTime(self::EPOCH, new \DateTimeZone('UTC'));
        $diff = $epoch->diff($date->setTimezone(new \DateTimeZone('UTC')));
        return $diff->days;
    }

    public static function dateFromDays(int $days): \DateTime {
        $epoch = new \DateTime(self::EPOCH, new \DateTimeZone('UTC'));
        return $epoch->add(new \DateInterval("P{$days}D"));
    }

    public static function secondsOfDay(\DateTimeInterface $time): int {
        return (int)$time->format('H') * 3600 + (int)$time->format('i') * 60 + (int)$time->format('s');
    }

    public static function timeFromSeconds(int $sec): \DateTime {
        $time = new \DateTime(self::EPOCH, new \DateTimeZone('UTC'));
        return $time->add(new \DateInterval("PT{$sec}S"));
    }

    public static function toUtcDateTime($o): \DateTime {
        if ($o instanceof \DateTimeInterface) {
            return $o->setTimezone(new \DateTimeZone('UTC'));
        }
        throw new \Exception('Unsupported datetime type');
    }

    public static function epochSeconds($o): int {
        if ($o instanceof \DateTimeInterface) {
            return (int)$o->setTimezone(new \DateTimeZone('UTC'))->format('U');
        }
        throw new \Exception('Unsupported datetime type');
    }
}
