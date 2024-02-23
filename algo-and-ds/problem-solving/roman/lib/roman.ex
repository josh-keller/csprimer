defmodule Roman do
  def convert(n, acc \\ [])

  def convert(0, acc) do
    acc
    |> Enum.reverse()
    |> Enum.join()
  end

  def convert(n, acc) when n > 0 do
    cond do
      n >= 1000 -> convert(n - 1000, ["M" | acc])
      n >= 900 -> convert(n - 900, ["CM" | acc])
      n >= 500 -> convert(n - 500, ["D" | acc])
      n >= 400 -> convert(n - 400, ["CD" | acc])
      n >= 100 -> convert(n - 100, ["C" | acc])
      n >= 90 -> convert(n - 90, ["XC" | acc])
      n >= 50 -> convert(n - 50, ["L" | acc])
      n >= 40 -> convert(n - 40, ["XL" | acc])
      n >= 10 -> convert(n - 10, ["X" | acc])
      n >= 9 -> convert(n - 9, ["IX" | acc])
      n >= 5 -> convert(n - 5, ["V" | acc])
      n >= 4 -> convert(n - 4, ["IV" | acc])
      n >= 1 -> convert(n - 1, ["I" | acc])
    end
  end
end
