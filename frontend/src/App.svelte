<script lang="ts">
    import logo from "./assets/images/logo-universal.png"
    import { Greet } from "../wailsjs/go/main/App.js"
    import { EventsOn } from "../wailsjs/runtime/runtime.js"

    let resultText: string = "Please enter your name below 👇"
    let name: string

    function greet(): void {
        Greet(name).then((result) => (resultText = result))
        getImage()
    }

    const getImage = async () => {
        try {
            const response = await fetch("screen.png")
            if (response.ok) {
                const blob = await response.blob()
                const imageUrl = URL.createObjectURL(blob)
                const imageElement = document.getElementById("screenshot") as HTMLImageElement
                imageElement.src = imageUrl
            } else {
                console.error("failed to fetch image")
            }
        } catch (error) {
            console.error("error:", error)
        }
    }

    EventsOn("capture", () => {
        greet()
    })
</script>

<main>
    <img alt="Sreentshot" id="screenshot" src={logo} />
</main>

<style>
    #screenshot {
        display: block;
        width: 100vw;
        height: 100vh;
    }
</style>
