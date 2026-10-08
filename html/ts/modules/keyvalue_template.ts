export class KeyValueTemplate {
    private static template_DOM: HTMLTemplateElement;

    public static SetTemplate(template: HTMLTemplateElement): void {
        KeyValueTemplate.template_DOM = template;
    }

    private key: string;
    private value: string;

    constructor(key: string, value: string) {
        this.key = key as string;
        this.value = value as string;
    }

    public Element(): HTMLTableRowElement {
        const tr_DOM = document.importNode(KeyValueTemplate.template_DOM.content, true).querySelector("tr") as HTMLTableRowElement;
        tr_DOM.children[0].innerHTML = this.key;
        tr_DOM.children[1].innerHTML = this.value;
        return tr_DOM as HTMLTableRowElement;
    }
}
