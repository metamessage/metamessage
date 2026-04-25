<?php

namespace io\metamessage\mm;

class GrowableByteBuf {
    private array $buf;
    private int $len;
    private const MAX_CAP = 1024 * 1024 * 1024; // 1GB

    public function __construct(int $initialSize = 1024) {
        $this->buf = array_fill(0, $initialSize, 0);
        $this->len = 0;
    }

    public function write(int ...$bs): void {
        $this->ensure(count($bs));
        foreach ($bs as $b) {
            $this->buf[$this->len++] = $b;
        }
    }

    public function writeAll(array $bs): void {
        $length = count($bs);
        if ($length === 0) {
            return;
        }
        $this->ensure($length);
        for ($i = 0; $i < $length; $i++) {
            $this->buf[$this->len + $i] = $bs[$i];
        }
        $this->len += $length;
    }

    public function size(): int {
        return $this->len;
    }

    public function copyRange(int $start, int $endExclusive): array {
        return array_slice($this->buf, $start, $endExclusive - $start);
    }

    public function reset(): void {
        $this->len = 0;
    }

    public function length(): int {
        return $this->len;
    }

    private function ensure(int $n): void {
        if ($this->len + $n > self::MAX_CAP) {
            throw new \Exception('Maximum size exceeded');
        }
        if ($this->len + $n > count($this->buf)) {
            $newCap = max(count($this->buf) * 2, $this->len + $n);
            $this->buf = array_pad($this->buf, $newCap, 0);
        }
    }
}
