import { KeyValueTemplate } from "./keyvalue_template";
export class ChoiceTemplate {
    private static template_DOM: HTMLTemplateElement;

    public static SetTemplate(template: HTMLTemplateElement): void {
        ChoiceTemplate.template_DOM = template;
    }

    private itemID: string = "";
    private tokenID: string;
    private dropdownVersion_DOM: HTMLElement;
    private modalVersion_DOM: HTMLElement;
    private modalButton_DOM: HTMLButtonElement;
    private dropdownCallback: ((choice: ChoiceTemplate) => void) | null = null;
    private modalCallback: ((choice: ChoiceTemplate) => void) | null = null;


    private elementBuilderHelper(imgSrc: string, data: KeyValueTemplate[], deleteButton: boolean): HTMLElement {
        const new_DOM = document.importNode(ChoiceTemplate.template_DOM.content, true).querySelector("aside") as HTMLElement;
        const img_DOM = new_DOM.querySelector("img") as HTMLImageElement;
        img_DOM.src = imgSrc;
        const table_DOM = new_DOM.querySelector("table") as HTMLTableElement;
        data.forEach(keyValue => {
            table_DOM.appendChild(keyValue.Element() as HTMLTableRowElement);
        });
        if (deleteButton) {
            const button = new_DOM.querySelector("button") as HTMLButtonElement;
            button.remove();
        }
        return new_DOM as HTMLElement;
    }

    constructor(tokenID: string, imgSrc: string, data: KeyValueTemplate[]) {
        this.tokenID = tokenID as string;
        this.dropdownVersion_DOM = this.elementBuilderHelper(imgSrc, data, true) as HTMLElement;
        this.modalVersion_DOM = this.elementBuilderHelper(imgSrc, data, false) as HTMLElement;
        this.modalButton_DOM = this.modalVersion_DOM.querySelector("button") as HTMLButtonElement;
    }

    private clickDropdownItem = (e: Event): void => {
        if (this.dropdownCallback) this.dropdownCallback(this);
    }

    private clickModalButton = (e: Event): void => {
        if (this.modalCallback) this.modalCallback(this);
    }

    public AddEventListeners(): void {
        this.dropdownVersion_DOM.addEventListener("click", this.clickDropdownItem);
        this.modalButton_DOM.addEventListener("click", this.clickModalButton);
    }

    public AddCallbacks(
        dropdownCallback: (choice: ChoiceTemplate) => void,
        modalCallback: (choice: ChoiceTemplate) => void,
    ): void {
        this.dropdownCallback = dropdownCallback;
        this.modalCallback = modalCallback;
    }

    public DropDownElement(): HTMLElement {
        return this.dropdownVersion_DOM as HTMLElement;
    }

    public ModalElement(): HTMLElement {
        return this.modalVersion_DOM as HTMLElement;
    }

    public SummaryElement(): HTMLElement {
        return this.dropdownVersion_DOM.cloneNode(true) as HTMLElement;
    }

    public set ItemID(itemID: string) {
        this.itemID = itemID;
    }

    public get TokenID(): string {
        return this.tokenID;
    }
}
