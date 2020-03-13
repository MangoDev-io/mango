import { Component, OnInit, Input } from '@angular/core'

@Component({
    selector: 'app-notification',
    templateUrl: './notification.component.html',
    styleUrls: ['./notification.component.scss'],
})
export class NotificationComponent implements OnInit {
    @Input()
    showModal = true

    @Input()
    modalSuccess = true

    @Input()
    assetId: string

    @Input()
    txHash: string

    @Input()
    error: string

    constructor() {}

    ngOnInit(): void {}

    toggleModal() {
        this.showModal = !this.showModal
    }

    getModalActiveClass() {
        if (this.showModal) {
            return 'is-active'
        }
    }

    shortenTxHash(): string {
        if (this.txHash)
            return (
                this.txHash.substring(0, 8) +
                ' . . . ' +
                this.txHash.substring(44)
            )
    }
}
