class LiveVideo {
    constructor(containerId, source) {
        this.m_containerId = containerId;
        this.m_container = document.getElementById(containerId);
        this.m_source = source;
    }

    display() {
        var video = document.getElementById("livevideo");
        if (video == null) {
            video = document.createElement("video");
            video.id = "livevideo";
            video.src = "http://192.168.0.94:3333";
            video.style.maxWidth = "100%";
            video.style.maxHeight = "100%";
            video.style.margin = "auto";
            this.m_container.appendChild(video);
        } else {
            video.style.visibility = "visible";
        }
    }

    hide() {
        var video = document.getElementById("livevideo");
        if (video != null) {
            video.style.visibility = "hidden";
        }
    }

    toString() {
        this.m_container.innerHTML = "LiveVideo";
    }
}

function ShowLiveVideo(containerId) {
    this.liveVideo = new LiveVideo(containerId);

    this.liveVideo.display();
}