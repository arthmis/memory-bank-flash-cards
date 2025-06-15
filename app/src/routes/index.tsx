import { getRouteApi, Link } from "@tanstack/solid-router";
import { createSignal, For } from "solid-js";

const fetchDecks = async (): Promise<Array<Deck>> => {
  const json = (await fetch("/api/decks")).json()
  return json
}

export const Route = createFileRoute({
  loader: async () => fetchDecks(),
  component: Index,
})

interface Deck {
  id: number;
  name: string;
}

function Index() {
  // TODO: add deck logic
  const [isAddingNewDeck, setIsAddingNewDeck] = createSignal(false);
  const [newDeckName, setNewDeckName] = createSignal("");
  const routeApi = getRouteApi('/')
  const data = routeApi.useLoaderData()
  const [decks, setDecks] = createSignal<Array<Deck>>(data());

  const addNewDeck = () => {
    fetch("/api/decks", {
      method: "POST",
      body: JSON.stringify({ name: newDeckName() }),
      headers: {
        "Content-Type": "application/json",
      },
    })
      .then((res) => res.json())
      .then((deck) => {
        setDecks([...decks(), deck]);
      });

    setIsAddingNewDeck(false);
    setNewDeckName("");
  }

  return (
    <div class="p-2">
      <h3>Welcome Home!</h3>
      <button onClick={() => setIsAddingNewDeck(true)}>New Deck</button>
      {isAddingNewDeck() && (
        <form>
          <input
            type="text"
            value={newDeckName()}
            onChange={(e) => setNewDeckName(e.target.value)}
          />
          <button onClick={addNewDeck}>Add Deck</button>
          <button onClick={() => setIsAddingNewDeck(false)}>Cancel</button>
        </form>
      )}
      <ul>
        <For each={decks()}>
          {(item, index) => (
            <li>
              <Link to={`/decks/${item.id}`}>{item.name}</Link>
            </li>
          )}
        </For>
      </ul>
    </div>
  )
}