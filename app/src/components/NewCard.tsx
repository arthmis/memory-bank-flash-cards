import { createSignal } from "solid-js";
import { useContext } from "solid-js"
import { JSX } from "solid-js/jsx-runtime";
import { Card, DeckContext, Route } from "../routes/decks/$deckId";

export function NewCard() {
    const { deckId } = Route.useParams()();
    console.log("deck id", deckId);

    const [question, setQuestion] = createSignal("");
    const [answer, setAnswer] = createSignal("");
    const {setState }  = useContext(DeckContext);

    const addCard: JSX.EventHandler<HTMLFormElement, SubmitEvent> = (event) => {
        event.preventDefault();

        const form = {
            question: question(),
            answer: answer()
        }

        fetch(`/api/decks/${deckId}/cards`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(form)
        }).then(res => res.json())
            .then(data => {
                console.log(data);
                setState("isAddingCard", false)
                setState("cards", (currentCards: Card[]) => [
                    ...currentCards,
                    data
                ])
            })
            .finally(() => {
                setQuestion("");
                setAnswer("");
            })
        console.log(event);
        // return the card and place into state
    }
    return (
        <>
            <form onSubmit={addCard} method="post">
                <label>Question</label>
                <input onInput={e => setQuestion(e.target.value)} value={question()} type="text" required />
                <label>Answer</label>
                <textarea onInput={e => setAnswer(e.target.value)} value={answer()} required />
                <button>Add</button>
            </form>
        </>
    )
}