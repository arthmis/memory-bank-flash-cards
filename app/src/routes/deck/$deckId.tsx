import { createContext, createSignal, For } from "solid-js";
import { createStore, Store } from "solid-js/store"
import { NewCard } from "../../components/NewCard";
import type { StoreNode, Store as StoreType, SetStoreFunction } from "solid-js/store"
import { getRouteApi } from "@tanstack/solid-router";

const fetchCards = async (): Promise<Array<Card>> => {
  const json = (await fetch("/api/deck/1/cards")).json()
  return json
}

export const Route = createFileRoute({
    // loader: async () => fetchCards(),
    component: DeckComponent,
})

interface DeckState {
    deckId: string;
    isAddingCard: boolean;
    cards: Card[]; // Define a proper type for cards if possible
}

export interface Card {
    question: string;
    answer: string;
}

export const DeckContext = createContext<{ state: Store<DeckState>, setState: SetStoreFunction<DeckState> }>();


function DeckComponent() {
    const { deckId } = Route.useParams()();
    // const routeApi = getRouteApi('/decks/$deckId/cards')
    // const cards = routeApi.useLoaderData()
    const [state, setState] = createStore({
        deckId,
        isAddingCard: false,
        // cards: cards(),
        cards: Array<Card>(),
    })
    const [isAddingCard, setIsAddingCard] = createSignal(false);
    return <>
        <DeckContext.Provider value={{ state, setState }}>
            <h1>{deckId}</h1>
            <button onClick={() => setIsAddingCard(true)}>New Card</button>
            {isAddingCard() && (
                <NewCard />
            )}
            <For each={state.cards}>
                {(item, index) => (
                    <div>
                        <p>{item.question}</p>
                        <p>{item.answer}</p>
                    </div>
                )}
            </For>
        </DeckContext.Provider>
    </>;
}