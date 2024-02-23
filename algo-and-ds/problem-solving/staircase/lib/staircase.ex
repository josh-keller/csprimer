defmodule Staircase do
  def climb(n) when n <= 0, do: 0
  def climb(n), do: _climb(n, 1, 1,0,0)

  def _climb(goal, curr, s1, _s2, _s3) when curr > goal, do: s1

  def _climb(goal, curr, s1, s2, s3) do
    _climb(goal, curr + 1, s1 + s2 + s3, s1, s2)
  end
end
