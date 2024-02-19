defmodule WordLadder do
  @wordlist File.read!("words.txt") |> String.upcase() |> String.split("\n", trim: true) |> MapSet.new()

  def process_words do
    result = case System.argv() do
      [_, word1, word2 | _rest] ->
        search(word1, word2)

      _ ->
        word1 = IO.gets("Enter the starting word: ") |> String.trim() |> String.upcase()
        word2 = IO.gets("Enter the goal word: ") |> String.trim() |> String.upcase()
        search(word1, word2)
    end
    
    result
    |> Enum.each(&IO.puts/1)
  end

  def search(start, target) do
    _search(target, [{[], [start]}], MapSet.new(), [], MapSet.new())
  end

  defp _search(target, [], visited, next_to_search, next_visited) do
    _search(target, next_to_search, MapSet.union(visited, next_visited), [], MapSet.new())
  end

  defp _search(target, [{history, words} | rest], visited, next_to_search, next_visited) do
    if Enum.member?(words, target) do
      [target | history] |> Enum.reverse()
    else
      new_neighbors = all_neighbors({history, words}, visited)
      _search(target, rest, visited, next_to_search ++ new_neighbors, MapSet.union(next_visited, visited_from_neigh(new_neighbors)))
    end
  end
  
  defp visited_from_neigh(neighbors) do
    neighbors
    |> Enum.flat_map(fn {_, words} -> words end)
    |> MapSet.new()
  end
  
  defp all_neighbors({history, words}, visited) do
    words
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

WordLadder.process_words()
