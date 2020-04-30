import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core'
import { StateService } from '../../state.service'

@Component({
    selector: 'app-notification',
    templateUrl: './notification.component.html',
    styleUrls: ['./notification.component.scss'],
})
export class NotificationComponent implements OnInit {
    @Input()
    showModal = true

    // 0 = success
    // 1 = failed
    // 2 = confirm
    @Input()
    notificationType = 0

    @Input()
    assetId: string

    @Input()
    txHash: string

    @Input()
    error: string

    @Output()
    confirmed = new EventEmitter<boolean>()

    @Output()
    modalClosed = new EventEmitter<boolean>()

    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    toggleModal() {
        this.showModal = !this.showModal
        this.modalClosed.emit(this.showModal)
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

    confirmDestroy() {
        this.confirmed.emit(true)
    }
}
