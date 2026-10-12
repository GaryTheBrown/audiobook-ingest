import { SendMessageType } from "../enums";
import { networkManager } from "../main";
export class SectionFailed {
    private itemID: string;
    private retryButton_DOM: HTMLButtonElement;
    private deleteButton_DOM: HTMLButtonElement;
    private message_DOM: HTMLParagraphElement;

    constructor(itemID: string, section_DOM: HTMLElement) {
        this.itemID = itemID;
        this.retryButton_DOM = section_DOM.querySelector("button.retry") as HTMLButtonElement;
        this.deleteButton_DOM = section_DOM.querySelector("button.delete") as HTMLButtonElement;
        this.message_DOM = section_DOM.querySelector("div p") as HTMLParagraphElement
    }

    private clickRetry = (e: Event): void => {
        const button = e.currentTarget as HTMLButtonElement;
        button.disabled = true
        networkManager.SendPayload(this.itemID, SendMessageType.FailedRetry, "");
    }

    private clickDelete = (e: Event): void => {
        const button = e.currentTarget as HTMLButtonElement;
        button.disabled = true
        networkManager.SendPayload(this.itemID, SendMessageType.FailedDelete, "");
    }

    public AddEventListeners(): void {
        this.retryButton_DOM.addEventListener("click", this.clickRetry)
        this.deleteButton_DOM.addEventListener("click", this.clickDelete)
    }

    public AddMessage(msg: string): void {
        this.message_DOM.innerText = msg;
    }
}
