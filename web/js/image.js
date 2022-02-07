class LiveImage {
    static sm_instance = null;

    constructor(containerId, source) {
        this.m_containerId = containerId;
        this.m_container = document.getElementById(containerId);
        this.m_source = source;
        if (LiveImage.sm_instance == null) {
            LiveImage.sm_instance = this;
        }
    }

    display() {
        var image = document.getElementById("liveimage");
        if (image == null) {
            image = document.createElement("img");
            image.id = "liveimage";
            image.src = "live";
            image.style.maxWidth = "100%";
            image.style.maxHeight = "100%";
            image.style.margin = "auto";
            this.m_container.appendChild(image);
        } else {
            image.style.visibility = "visibile";
        }
    }

    hide() {
        var image = document.getElementById("liveimage");
        if (image != null) {
            image.style.visibility = "hidden";
        }
    }

    toString() {
        this.m_container.innerHTML = "LiveImage";
    }
}

function ShowLiveImage(containerId) {
    if (LiveImage.sm_instance == null) {
        this.liveImage = new LiveImage(containerId);
    }
    LiveImage.sm_instance.display();
}

function HideLiveImage() {
    LiveImage.sm_instance.hide();
}