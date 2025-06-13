export const Route = createFileRoute({
  component: Index,
})

function Index() {
  return (
    <div class="p-2">
      <h3>Welcome Home!</h3>
      <button>New Deck</button>
    </div>
  )
}