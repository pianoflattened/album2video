module.exports = class ProgressBar {
	constructor(DOMelement) {
		this._DOMelement = DOMelement;
		this._bar = DOMelement.querySelector(".progress-bar");
		this._label = DOMelement.querySelector("#progress-label");
	}
	
	setProgress(progress) {
		this._bar.setAttribute("aria-valuenow", progress.toString());
	}
	
	setLabel(text) {
		this._label.textContent = text;
	}
	
	makeDeterminate() {
		this._bar.setAttribute("aria-valuenow", "0");
		this._bar.setAttribute("aria-valuemin", "0");
		this._bar.setAttribute("aria-valuemax", "100");
		this._bar.class = "progress-bar";
	}
	
	makeIndeterminate() {
		this._bar.setAttribute("aria-valuenow", "100");
		this._bar.setAttribute("aria-valuemin", "0");
		this._bar.setAttribute("aria-valuemax", "100");
		this._bar.class = "progress-bar progress-bar-striped progress-bar-animated";
	}
	
	setComplete() {
		this._bar.setAttribute("aria-valuenow", "100");
		if (this._bar.classList.contains("progress-bar-animated")) {
			this._bar.classList.remove("progress-bar-animated");
		}
		this.setLabel("done!");
	}
	
	error(msg) {
		this._bar.setAttribute("aria-valuenow", "100");
		if (this._bar.classList.contains("progress-bar-animated")) {
			this._bar.classList.remove("progress-bar-animated");
		}
		this.setLabel(msg);
	}
}