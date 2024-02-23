defmodule WordLadder do
  # Read in the word list to a global Set
  @wordlist File.read!("words.txt") |> String.upcase() |> String.split("\n", trim: true) |> MapSet.new()

  def process_words do
    # Get words from command line or interactively
    {start, target} = case System.argv() do
      [word1, word2] ->
        {String.upcase(word1), String.upcase(word2)}

      [] ->
        word1 = IO.gets("Enter the starting word: ") |> String.trim() |> String.upcase()
        word2 = IO.gets("Enter the goal word: ") |> String.trim() |> String.upcase()
        {word1, word2}

      _ ->
        IO.inspect(System.argv())
        IO.puts("Usage: pass start and end words as command line arguments, or pass no args to run interactively")
        System.stop(1)
    end

    # Check words are valid
    if not MapSet.member?(@wordlist, start) do
      IO.puts("#{start} not in word list")
      System.stop(1)
    end
    if not MapSet.member?(@wordlist, target) do
      IO.puts("#{target} not in word list")
      System.stop(1)
    end
    if String.length(start) != String.length(target) do
      IO.puts("Words must be same length")
      System.stop(1)
    end
    
    if result = search(start, target) do
      result
      |> Enum.each(&IO.puts/1)
    else
      IO.puts("No word ladder found")
    end
  end

  def search(start, target) do
    _search(target, [{[], [start]}], MapSet.new(), [], MapSet.new())
  end

  # If the current search level and next search level are empty, the search has failed
  defp _search(_target, [], _visited, [], _next_visited), do: nil

  # If the current search level is empty, it is time to search the next one
  defp _search(target, [], visited, next_to_search, next_visited) do
    _search(target, next_to_search, MapSet.union(visited, next_visited), [], MapSet.new())
  end

  # Iterate through the current search level
  defp _search(target, [{history, words} | rest], visited, next_to_search, next_visited) do
    if Enum.member?(words, target) do
      [target | history] |> Enum.reverse()
    else
      new_neighbors = all_neighbors({history, words}, visited)
      _search(
        target,
        rest,
        visited,
        next_to_search ++ new_neighbors,
        MapSet.union(next_visited, extract_visited(new_neighbors))
      )
    end
  end
  
  defp extract_visited(neighbors) do
    neighbors
    |> Enum.flat_map(fn {_, words} -> words end)
    |> MapSet.new()
  end
  
  defp all_neighbors({history, words}, visited) do
    words
    |> Enum.map(fn word -> 
      {
        [word | history],
        neighbors(word)
        |> Enum.filter(&(!MapSet.member?(visited, &1)))
      }
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
