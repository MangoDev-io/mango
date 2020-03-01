import { Injectable } from '@angular/core'
import { Token } from './model/token'
import { BehaviorSubject, Observable } from 'rxjs'

@Injectable({
    providedIn: 'root',
})
export class StateService {
    public currToken = new Token(
        '324234',
        '2LX7ZMR7SMDONF3FLD2SM5KUSKUWYKDH4WS76AW26US3Y3QB4Z4UROVFTY',
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

    private selectedTokenSubject = new BehaviorSubject<Token>(null)
    private showCreateSubject = new BehaviorSubject<boolean>(false)

    constructor() {
        this.showCreateSubject.next(true)
    }

    getMnemonic(): string {
        return this.mnemonicPhrase
    }

    setMnemonic(mnemonic: string) {
        this.mnemonicPhrase = mnemonic
    }

    getSelectedToken(): Observable<Token> {
        return this.selectedTokenSubject.asObservable()
    }

    setSelectedToken(token: Token) {
        this.selectedTokenSubject.next(token)
    }

    getShowCreateToken(): Observable<boolean> {
        return this.showCreateSubject.asObservable()
    }

    setShowCreateToken(b: boolean) {
        this.showCreateSubject.next(b)
    }
}
