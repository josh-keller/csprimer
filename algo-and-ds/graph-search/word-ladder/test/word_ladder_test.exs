defmodule WordLadderTest do
  use ExUnit.Case
  doctest WordLadder

  test "find neighbors" do
    assert WordLadder.neighbors("WORD") ==["BORD", "CORD", "FORD", "LORD", "OORD", "WARD", "WIRD", "WOAD", "WOLD", "WOOD", "WORE", "WORK", "WORM", "WORN", "WORT"] 
  end
end
