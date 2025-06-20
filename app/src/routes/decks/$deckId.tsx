import { createContext, createSignal, For } from "solid-js";
import { createStore, Store } from "solid-js/store"
import { NewCard } from "../../components/NewCard";
import type { StoreNode, Store as StoreType, SetStoreFunction } from "solid-js/store"
import { getRouteApi } from "@tanstack/solid-router";
import { Card, getCardsByDeckId, getCardsByDeckIdResponse } from "../../orval-client";

const fetchCards = async (params): Promise<getCardsByDeckIdResponse> => {
    const { deckId } = params;
    const data = getCardsByDeckId(deckId)
    return data
}

export const Route = createFileRoute({
    loader: async ({ params }) => fetchCards(params),
    component: DeckComponent,
})

interface DeckState {
    deckId: string;
    isAddingCard: boolean;
    cards: Card[];
}

export const DeckContext = createContext<{ state: Store<DeckState>, setState: SetStoreFunction<DeckState> }>();


function DeckComponent() {
    const { deckId } = Route.useParams()();
    const { data } = Route.useLoaderData()()
    const [state, setState] = createStore({
        deckId,
        isAddingCard: false,
        cards: data.cards ?? [],
    })

    return <>
        <DeckContext.Provider value={{ state, setState }}>
            <h1>{deckId}</h1>
            <button onClick={() => setState("isAddingCard", true)}>New Card</button>
            {state.isAddingCard && (
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