const sleep = require('./sleep.js'); 

module.exports = class ProgressBar {
	constructor(DOMelement) {
		this._DOMelement = DOMelement;
		this._bar = DOMelement.querySelector(".progress-bar");
		this._label = DOMelement.querySelector("#progress-label");
	}
	
	async setProgress(progress) {
		this._bar.setAttribute("aria-valuenow", progress.toString());
		this._bar.setAttribute("style", "width: " + progress + "%");
        await sleep(100);
	}
	
	async setLabel(text) {
		this._label.textContent = text;
        await sleep(100);
	}
	
	async makeDeterminate() {
		this._bar.setAttribute("aria-valuenow", "0");
		this._bar.setAttribute("aria-valuemin", "0");
		this._bar.setAttribute("aria-valuemax", "100");
		this._bar.setAttribute("style", "width: 0%");
		this._bar.class = "progress-bar";
        await sleep(100);
	}
	
	async makeIndeterminate() {
		this._bar.setAttribute("aria-valuenow", "100");
		this._bar.setAttribute("aria-valuemin", "0");
		this._bar.setAttribute("aria-valuemax", "100");
		this._bar.setAttribute("style", "width: 100%");
		this._bar.setAttribute("class", "progress-bar progress-bar-striped progress-bar-animated");
        await sleep(100);
	}
	
	async setComplete() {
		this._bar.setAttribute("aria-valuenow", "100");
		this._bar.setAttribute("style", "width: 100%");
		if (this._bar.classList.contains("progress-bar-animated")) {
			this._bar.classList.remove("progress-bar-animated");
		}
        if (this._bar.classList.contains("progress-bar-striped")) {
			this._bar.classList.remove("progress-bar-striped");
		}
		this.setLabel("done!");
        await sleep(100);
	}
	
	async error(msg) {
		this._bar.setAttribute("aria-valuenow", "100");
		this._bar.setAttribute("style", "width: 100%");
		if (this._bar.classList.contains("progress-bar-animated")) {
			this._bar.classList.remove("progress-bar-animated");
		}
		this.setLabel(msg);
        await sleep(100);
	}
}
