import { KeyValueTemplate } from "../modules/keyvalue_template";

export class SectionConverting {
    private asideImg_DOM: HTMLImageElement;
    private asideTable_DOM: HTMLTableElement;
    private progressBar_DOM: HTMLElement;
    private progressBarSpan_DOM: HTMLSpanElement;
    private progressInfo_DOM: HTMLElement;

    constructor(section_DOM: HTMLElement) {
        this.asideImg_DOM = section_DOM.querySelector("aside img") as HTMLImageElement;
        this.asideTable_DOM = section_DOM.querySelector("aside table") as HTMLTableElement;
        this.progressBar_DOM = section_DOM.querySelector("div.progress") as HTMLTableElement;
        this.progressBarSpan_DOM = this.progressBar_DOM.querySelector("span") as HTMLTableElement;
        this.progressInfo_DOM = section_DOM.querySelector("div.info") as HTMLTableElement;
    }

    public SetImage(imgSrc: string): void {
        this.asideImg_DOM.src = imgSrc
    }

    public AddMetadata(keyValues: KeyValueTemplate[]): void {
        this.asideTable_DOM.innerHTML = "";
        keyValues.forEach(keyValue => {
            this.asideTable_DOM.appendChild(keyValue.Element() as HTMLTableRowElement);
        });
    }

    public SetProgress(percent: number, text: string): void {
        const percentageString = `${percent.toFixed(0)}%`;
        this.progressBar_DOM.setAttribute('data-label', percentageString);
        this.progressBarSpan_DOM.style.width = percentageString;
        this.progressInfo_DOM.innerText = text;
    }
}