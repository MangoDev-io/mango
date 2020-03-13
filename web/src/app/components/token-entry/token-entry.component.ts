import { Component, OnInit, Input, ElementRef } from '@angular/core'
import { StateService } from 'src/app/state.service'
import { Token } from 'src/app/token'
import { ThrowStmt } from '@angular/compiler'

@Component({
    selector: 'app-token-entry',
    templateUrl: './token-entry.component.html',
    styleUrls: ['./token-entry.component.scss'],
})
export class TokenEntryComponent implements OnInit {
    @Input()
    assetId: string

    @Input()
    permissions: string[]

    @Input()
    entryId: any

    token: Token

    gradient: string

    constructor(private stateService: StateService) {}

    ngOnInit() {
        this.getAssetDetails(this.assetId)
        this.generateGradient()
    }

    private generateGradient(): void {
        // prettier-ignore
        var hexValues = ["0","1","2","3","4","5","6","7","8","9","a","b","c","d","e"];

        function populate(a) {
            for (var i = 0; i < 6; i++) {
                var x = Math.round(Math.random() * 14)
                var y = hexValues[x]
                a += y
            }
            return a
        }

        var newColor1 = populate('#')
        var newColor2 = populate('#')
        var angle = Math.round(Math.random() * 360)

        // prettier-ignore
        this.gradient = "linear-gradient(" + angle + "deg, " + newColor1 + ", " + newColor2 + ")";
    }

    changeSelectedToken(): void {
        this.stateService.setSelectedToken(this.token)
        this.stateService.setShowCreateToken(false)
        let entries = document.getElementsByClassName('token-entry__container')
        for (let i = 0; i < entries.length; i++) {
            entries[i].classList.remove('active')
        }

        let createNewEntries = document.getElementsByClassName(
            'create-new__container'
        )
        for (let i = 0; i < createNewEntries.length; i++) {
            createNewEntries[i].classList.remove('active')
        }

        let curEntry = document.getElementById(`tokenEntry__${this.entryId}`)
        curEntry.firstElementChild.classList.add('active')
    }

    async getAssetDetails(assetId: string) {
        let assetInfo = await this.stateService.getAssetDetails(assetId)
        this.token = new Token({
            assetId: assetId,
            creatorAddr: assetInfo.creator,
            assetName: assetInfo.assetname,
            unitName: assetInfo.unitname,
            totalIssuance: assetInfo.total,
            decimals: assetInfo.decimals,
            defaultFrozen: assetInfo.defaultfrozen,
            url: assetInfo.url,
            metadataHash: assetInfo.metadatahash,
            managerAddr: assetInfo.managerkey,
            reserveAddr: assetInfo.reserveaddr,
            freezeAddr: assetInfo.freezeaddr,
            clawbackAddr: assetInfo.clawbackaddr,
            permissions: this.permissions,
        })
    }
}
