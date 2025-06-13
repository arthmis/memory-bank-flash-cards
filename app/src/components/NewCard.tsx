import { createSignal } from "solid-js";
import { useContext } from "solid-js"
import { JSX } from "solid-js/jsx-runtime";

export function NewCard() {
    const [question, setQuestion] = createSignal("");
    const [answer, setAnswer] = createSignal("");
    const addCard: JSX.EventHandler<HTMLFormElement, SubmitEvent> = (event) => {
        event.preventDefault();

        const form = {
            question: question(),
            answer: answer()
        }

        fetch("/api/cards", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(form)
        }).then(res => res.json())
            .then(data => {
                console.log(data);
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