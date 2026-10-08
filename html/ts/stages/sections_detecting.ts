import { ListItemTemplate } from "../modules/list_item_template";
export class SectionDetecting {

    private ol_DOM: HTMLOListElement;
    constructor(section_DOM: HTMLElement) {
        this.ol_DOM = section_DOM.querySelector("ol") as HTMLOListElement;
    }

    public AddLog(ListItem: ListItemTemplate): void {
        this.ol_DOM.appendChild(ListItem.Element());
    }

    public RemoveLogs(): void {
        const elementsToRemove = this.ol_DOM.querySelectorAll('li:not(:first-child)');
        for (const li of elementsToRemove) {
            li.remove();
        }
    }
}
