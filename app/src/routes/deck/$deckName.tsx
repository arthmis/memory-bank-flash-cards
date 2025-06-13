import { createSignal } from "solid-js";
import { NewCard } from "../../components/NewCard";

export const Route = createFileRoute({
    component: DeckComponent,
})

function DeckComponent() {
    const { deckName } = Route.useParams()();
    const [isAddingCard, setIsAddingCard] = createSignal(false);
    return <>
        <h1>{deckName}</h1>
        <button onClick={() => setIsAddingCard(true)}>New Card</button>
        {isAddingCard() && (
            <NewCard />
        )}
    </>
}