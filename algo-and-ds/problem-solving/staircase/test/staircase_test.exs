defmodule StaircaseTest do
  use ExUnit.Case
  doctest Staircase

  test "staircase ascent" do
    assert Staircase.climb(1) == 1
    assert Staircase.climb(2) == 2
    assert Staircase.climb(3) == 4
    assert Staircase.climb(4) == 7
    assert Staircase.climb(5) == 13
    assert Staircase.climb(20) == 121415
  end
end
