import { createSignal } from "solid-js";
import { useContext } from "solid-js"
import { JSX } from "solid-js/jsx-runtime";
import { DeckContext, Route } from "../routes/decks/$deckId";
import { Card, createCard } from "../orval-client";

export function NewCard() {
    const { deckId } = Route.useParams()();
    console.log("deck id", deckId);

    const [question, setQuestion] = createSignal("");
    const [answer, setAnswer] = createSignal("");
    const { setState, } = useContext(DeckContext);

    const addCard: JSX.EventHandler<HTMLFormElement, SubmitEvent> = (event) => {
        event.preventDefault();

        const form = {
            question: question(),
            answer: answer()
        }

        createCard(Number.parseInt(deckId), form
        )
            .then((res) => {
                setState("isAddingCard", false)
                setState("cards", (currentCards: Card[]) => [
                    ...currentCards,
                    res.data
                ])
            })
            .finally(() => {
                setQuestion("");
                setAnswer("");
            })
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