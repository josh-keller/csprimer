defmodule RomanTest do
  use ExUnit.Case
  doctest Roman

  test "converts int to roman numeral" do
    test_cases = %{
      39 => "XXXIX",
      246 => "CCXLVI",
      789 => "DCCLXXXIX",
      2421 => "MMCDXXI",
      160 => "CLX",
      207 => "CCVII",
      1009 => "MIX",
      1066 => "MLXVI",
      1776 => "MDCCLXXVI",
      1918 => "MCMXVIII",
      1944 => "MCMXLIV",
      2024 => "MMXXIV"
    }
    
    test_cases
    |> Enum.each(fn {n, r} -> assert(Roman.convert(n) == r) end)
  end
end
