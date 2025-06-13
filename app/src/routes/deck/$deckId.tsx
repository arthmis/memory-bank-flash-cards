import { createContext, createSignal } from "solid-js";
import { createStore, Store } from "solid-js/store"
import { NewCard } from "../../components/NewCard";
import type { StoreNode, Store as StoreType, SetStoreFunction } from "solid-js/store"

export const Route = createFileRoute({
    component: DeckComponent,
})

interface DeckState {
    deckId: string;
    cards: Card[]; // Define a proper type for cards if possible
}

interface Card {
    question: string;
    answer: string;
}

const DeckContext = createContext<{ state: Store<DeckState>, setState: SetStoreFunction<DeckState> }>();


function DeckComponent() {
    const { deckId } = Route.useParams()();
    const [state, setState] = createStore({
        deckId,
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
        </DeckContext.Provider>
    </>;
}