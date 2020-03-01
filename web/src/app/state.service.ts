import { Injectable } from '@angular/core'
import { Token } from './model/token'

@Injectable({
    providedIn: 'root',
})
export class StateService {
    public currToken = new Token(
        '324234',
        'JHJBFDISFD43534534534FGDSGF',
        'USD Tether',
        'USDT',
        100000000,
        15,
        false,
        'https://usdtether.io',
        '',
        '',
        '',
        '',
        ''
    )

    public tokenList: Token[] = [this.currToken, this.currToken, this.currToken]

    private mnemonicPhrase: string

    constructor() {}

    getMnemonic(): string {
        return this.mnemonicPhrase
    }

    setMnemonic(mnemonic: string) {
        this.mnemonicPhrase = mnemonic
    }
}
