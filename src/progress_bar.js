const ElectronProgressBar = require('electron-progressbar');

class ProgressBar extends ElectronProgressBar {
	constructor(options, electronApp) {
		super(options, electronApp);
	}
	
	error(text) {
		this.text = "ERROR";
		this.detail = text;
		this._realValue = this._options.maxValue;
		this._updateTaskbarProgress();
		if (this._options.indeterminate) {
			this._options.style.bar.background = this._options.style.value.background;
		}
	}
}

module.exports = ProgressBar;