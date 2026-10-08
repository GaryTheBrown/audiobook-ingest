import { ListItemTemplate } from "../modules/list_item_template";
export class SectionDetecting {
    private itemID: string;

    private ol_DOM: HTMLOListElement;
    constructor(itemID: string, section_DOM: HTMLElement) {
        this.itemID = itemID;
        this.ol_DOM = section_DOM.querySelector("ol") as HTMLOListElement;
    }
    public AddEventListeners(): void { }

    public AddLog(ListItem: ListItemTemplate) {
        this.ol_DOM.appendChild(ListItem.Element());
    }

    public RemoveLogs() {
        const elementsToRemove = this.ol_DOM.querySelectorAll('li:not(:first-child)');
        for (const li of elementsToRemove) {
            li.remove();
        }
    }
}
