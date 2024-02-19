defmodule WordLadder do
  @wordlist File.read!("words.txt") |> String.upcase() |> String.split("\n", trim: true) |> MapSet.new()

  def search(start, target) do
    _search(target, [{[], [start]}], MapSet.new(), [], MapSet.new())
  end

  # defp _search(target, to_search, visited, next_to_search, next_visited)
  defp _search(target, [], visited, next_to_search, next_visited) do
    _search(target, next_to_search, MapSet.union(visited, next_visited), [], MapSet.new())
  end

  defp _search(target, [{history, words} | rest], visited, next_to_search, next_visited) do
    IO.inspect(history)
    IO.inspect(words)
    IO.inspect(rest)
    IO.puts("-------")
    if Enum.member?(words, target) do
      [target | history] |> Enum.reverse()
    else
      new_neighbors = all_neighbors({history, words}, visited)
      IO.inspect(new_neighbors)
      _search(target, rest, visited, next_to_search ++ new_neighbors, MapSet.union(next_visited, visited_from_neigh(new_neighbors)))
    end
  end
  
  defp visited_from_neigh(neighbors) do
    neighbors
    |> Enum.flat_map(fn {_, words} -> words end)
    |> MapSet.new()
  end
  
  defp all_neighbors({history, words}, visited) do
    IO.puts("alln")
    words
    |> IO.inspect
    |> Enum.map(fn word -> 
      {[word | history], neighbors(word) |> Enum.filter(&(!MapSet.member?(visited, &1)))}
    end)
  end

  def neighbors(word) do
    _neighbors("", word, [])
  end

  defp _neighbors(_, "", neighbors), do: neighbors

  defp _neighbors(pre, <<c>> <> rest, neighbors) do
    next_neighbors = ~c(ABCDEFGHIJKLMNOPQRSTUVWXYZ)
    |> Enum.filter(fn char -> char != c end)
    |> Enum.map(fn char -> pre <> <<char>> <> rest end)
    |> Enum.filter(fn word -> MapSet.member?(@wordlist, word) end)

    _neighbors(pre <> <<c>>, rest, neighbors ++ next_neighbors)
  end
end

WordLadder.search("HEAD", "TAIL") |> IO.inspect
