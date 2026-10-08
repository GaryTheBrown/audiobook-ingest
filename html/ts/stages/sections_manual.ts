import { SendMessageType } from "../enums";
import { networkManager } from "../main";

export class SectionManual {
    private itemID: string;
    private section_DOM: HTMLElement;
    private btnSearch: HTMLButtonElement;
    private btnReScan: HTMLButtonElement;

    constructor(itemID: string, section_DOM: HTMLElement) {
        this.itemID = itemID;
        this.section_DOM = section_DOM;
        this.btnSearch = this.section_DOM.querySelector('footer button:first-child') as HTMLButtonElement;
        this.btnReScan = this.section_DOM.querySelector('footer button:last-child') as HTMLButtonElement;
    }

    public AddEventListeners(): void {
        if (this.btnSearch) this.btnSearch.addEventListener("click", this.ButtonClickSearch);
        if (this.btnReScan) this.btnReScan.addEventListener("click", this.ButtonClickReScan);
    }

    private ButtonClickSearch = (e: Event): void => {
        if (this.btnSearch) this.btnSearch.setAttribute("disabled", "true");
        networkManager.SendPayload(this.itemID, SendMessageType.ManualSearch, this.GetInputs());
    }

    private ButtonClickReScan = (e: Event): void => {
        const target = e.currentTarget as HTMLElement;
        if (target) target.setAttribute("disabled", "true");
        networkManager.SendPayload(this.itemID, SendMessageType.Rescan, "");
    }

    private GetInputs(): string {
        const inputList = this.section_DOM.querySelectorAll("input") as NodeListOf<HTMLInputElement>;
        const Data = new Map<string, string>();
        inputList.forEach((input: HTMLInputElement) => {
            const inputName = input.getAttribute("name");
            if (inputName) Data.set(inputName, input.value);
        });
        return JSON.stringify(Object.fromEntries(Data));
    }
}
