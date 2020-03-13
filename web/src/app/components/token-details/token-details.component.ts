import { Component, OnInit, Input } from '@angular/core'
import { Token } from '../../token'
import { AssetRequest } from '../../assetRequest'
import { StateService } from 'src/app/state.service'

@Component({
    selector: 'app-token-details',
    templateUrl: './token-details.component.html',
    styleUrls: ['./token-details.component.scss'],
})
export class TokenDetailsComponent implements OnInit {
    @Input()
    currToken: Token

    createButtonLoading = false

    selectedButton = 1
    assetManagementTabs = {
        '1': 'Freeze',
        '2': 'Revoke',
        '3': 'Modify',
        '4': 'Destroy',
    }

    constructor(private stateService: StateService) {}

    ngOnInit(): void {
        console.log('Token')
        console.log(this.currToken)
    }

    shortenAddress(addr: string): string {
        if (addr) return addr.substring(0, 8) + ' . . . ' + addr.substring(50)
    }

    getButtonLoadingClass() {
        if (this.createButtonLoading) {
            return 'is-loading'
        }
    }

    updateSelectedButton(b: number) {
        this.selectedButton = b
        console.log(
            'Selected tab: ' + this.assetManagementTabs[this.selectedButton]
        )
    }

    assetRequest = new AssetRequest()
    handleAssetRequest() {
        this.assetRequest.assetId = parseInt(this.currToken.assetId)
        this.createButtonLoading = true
        switch (this.selectedButton) {
            case 1: {
                console.log(
                    'Freeze request: ' + JSON.stringify(this.assetRequest)
                )
                this.stateService
                    .freezeAsset(this.assetRequest)
                    .subscribe(x => {
                        console.log(x)
                        this.createButtonLoading = false
                    })
                break
            }

            case 2: {
                console.log(
                    'Revoke request: ' + JSON.stringify(this.assetRequest)
                )
                this.stateService
                    .revokeAsset(this.assetRequest)
                    .subscribe(x => {
                        console.log(x)
                        this.createButtonLoading = false
                    })
                break
            }

            case 3: {
                console.log(
                    'Modify request: ' + JSON.stringify(this.assetRequest)
                )
                this.stateService
                    .modifyAsset(this.assetRequest)
                    .subscribe(x => {
                        console.log(x)
                        this.createButtonLoading = false
                    })
                break
            }

            case 4: {
                console.log(
                    'Destroy request: ' + JSON.stringify(this.assetRequest)
                )
                this.stateService
                    .destroyAsset(this.assetRequest)
                    .subscribe(x => {
                        console.log(x)
                        this.createButtonLoading = false
                    })
                break
            }
        }
    }
}
