<script lang="ts">
    import logo from "./assets/images/logo-universal.png"
    import { EventsOn } from "../wailsjs/runtime/runtime.js"

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
        getImage()
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
