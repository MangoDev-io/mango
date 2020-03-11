import { Component, OnInit, Input } from '@angular/core'
import { StateService } from 'src/app/state.service'
import { Token } from '../../token'
import { AssetRequest } from '../../assetRequest'

@Component({
    selector: 'app-token-details',
    templateUrl: './token-details.component.html',
    styleUrls: ['./token-details.component.scss'],
})
export class TokenDetailsComponent implements OnInit {
    @Input()
    currToken: Token

    selectedButton = 1

    constructor(private stateService: StateService) {}

    ngOnInit(): void {}

    shortenAddress(addr: string): string {
        if (addr) return addr.substring(0, 8) + ' . . . ' + addr.substring(50)
    }

    updateSelectedButton(b: number) {
        this.selectedButton = b
        console.log(this.selectedButton)
    }

    assetRequest = new AssetRequest(); 
    handleAssetRequest() {
        switch (this.selectedButton) {
            case 1: {
                console.log("Freeze request: " + JSON.stringify(this.assetRequest))
                this.stateService.freezeAsset(this.assetRequest).subscribe(x => {
                    console.log(x)
                })
                break;
            }
            
            case 2: {
                console.log("Revoke request: " + JSON.stringify(this.assetRequest))
                this.stateService.revokeAsset(this.assetRequest).subscribe(x => {
                    console.log(x)
                })
                break;
            }

            case 3: {
                console.log("Modify request: " + JSON.stringify(this.assetRequest))
                this.stateService.modifyAsset(this.assetRequest).subscribe(x => {
                    console.log(x)
                })
                break;
            }

            case 4: {
                console.log("Destroy request: " + JSON.stringify(this.assetRequest))
                this.stateService.destroyAsset(this.assetRequest).subscribe(x => {
                    console.log(x)
                })
                break;
            }
        }
    }
}
