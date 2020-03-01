import { Component, OnInit } from '@angular/core'
import { Token } from '../../model/token'

@Component({
    selector: 'app-token-details',
    templateUrl: './token-details.component.html',
    styleUrls: ['./token-details.component.scss'],
})
export class TokenDetailsComponent implements OnInit {
    public currToken = new Token(
        '324234',
        'BVMEUTF37WNEQ6GYCZISRFHGLEMOKT5OCPPTTJXVED6JBSXKF6YJJRZRI4',
        'USD Tether',
        'USDT',
        100000000,
        15,
        false,
        'https://tether.to',
        'sdfsdf',
        '',
        '',
        '',
        ''
    )

    constructor() {}

    ngOnInit(): void {}

    shortenAddress(addr: string): string {
        if (addr) return addr.substring(0, 10) + '. . .' + addr.substring(48)
    }
}
