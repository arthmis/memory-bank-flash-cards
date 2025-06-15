import { createContext, createSignal, For } from "solid-js";
import { createStore, Store } from "solid-js/store"
import { NewCard } from "../../components/NewCard";
import type { StoreNode, Store as StoreType, SetStoreFunction } from "solid-js/store"
import { getRouteApi } from "@tanstack/solid-router";

const fetchCards = async ( params ): Promise<Array<Card>> => {
  const { deckId } = params;
  const json = (await fetch(`/api/decks/${deckId}/cards`)).json()
  return json
}

export const Route = createFileRoute({
    // loaderDeps: ({ params }) > ({ params }),
    loader: async ({ params}) => fetchCards(  params  ),
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
    const cards = Route.useLoaderData()
    const [state, setState] = createStore({
        deckId,
        isAddingCard: false,
        cards: cards(),
        // cards: Array<Card>(),
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