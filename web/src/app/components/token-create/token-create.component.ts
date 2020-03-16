import { Component, OnInit } from '@angular/core'
import { StateService } from 'src/app/state.service'
import { Token } from '../../token'

@Component({
    selector: 'app-token-create',
    templateUrl: './token-create.component.html',
    styleUrls: ['./token-create.component.scss'],
})
export class TokenCreateComponent implements OnInit {
    assetCreate = new Token()

    createButtonLoading = false

    showNotificationModal = false
    notificationModalSuccess = true

    responseAssetId = ''
    responseTxHash = ''
    responseError = ''

    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    createAsset() {
        console.log(JSON.stringify(this.assetCreate))
        this.createButtonLoading = true
        this.stateService.createAsset(this.assetCreate).subscribe(
            x => {
                console.log(x)
                this.createButtonLoading = false
                this.showNotificationModal = true
                this.clearForm()
                this.responseAssetId = x.assetId.toString()
                this.responseTxHash = x.txHash
                this.stateService.setReloadListings()
            },
            err => {
                console.error(err)
                this.createButtonLoading = false
                this.showNotificationModal = true
                this.notificationModalSuccess = false
                this.clearForm()
                this.responseError = err.error.message
            }
        )
    }

    clearForm() {
        this.assetCreate = new Token()
    }

    getButtonLoadingClass() {
        if (this.createButtonLoading) {
            return 'is-loading'
        }
    }
}
