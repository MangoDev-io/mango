import { Component, OnInit } from '@angular/core'

@Component({
    selector: 'app-notification',
    templateUrl: './notification.component.html',
    styleUrls: ['./notification.component.scss'],
})
export class NotificationComponent implements OnInit {
    showModal = true
    modalSuccess = true

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

    shortenTxHash(hash: string): string {
        if (hash) return hash.substring(0, 8) + ' . . . ' + hash.substring(44)
    }
}
