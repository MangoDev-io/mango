import { Component, OnInit, Input } from '@angular/core'
import { Token } from '../../token'

@Component({
    selector: 'app-token-details',
    templateUrl: './token-details.component.html',
    styleUrls: ['./token-details.component.scss'],
})
export class TokenDetailsComponent implements OnInit {
    @Input()
    currToken: Token

    constructor() {}

    ngOnInit(): void {}

    shortenAddress(addr: string): string {
        if (addr) return addr.substring(0, 8) + ' . . . ' + addr.substring(50)
    }
}
