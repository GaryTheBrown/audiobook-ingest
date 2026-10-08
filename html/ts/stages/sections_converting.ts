import { KeyValueTemplate } from "../modules/keyvalue_template";

export class SectionConverting {
    private consoleLog_DOM: HTMLDialogElement;
    private asideImg_DOM: HTMLImageElement;
    private asideTable_DOM: HTMLTableElement;

    constructor(section_DOM: HTMLElement) {
        this.consoleLog_DOM = section_DOM.querySelector("details dialog") as HTMLDialogElement;
        this.asideImg_DOM = section_DOM.querySelector("aside img") as HTMLImageElement;
        this.asideTable_DOM = section_DOM.querySelector("aside table") as HTMLTableElement;
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
}