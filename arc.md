                         USER
                          │
                    ┌─────┴─────┐
                    │           │
                  WALLET      AGENT
                    │           │
                    │      ┌────┴────┐
                    │      │         │
                    │    POLICY    SESSION
                    │      │         │
                    └──────┴─────────┘
                             │
                         PAYMENT
                             │
                           x402
                             │
                       API SERVICE
                             │
                       POLICY ENGINE
                             │
                     ┌───────┴───────┐
                     │               │
                   BLOCK           APPROVE
                     │               │
                  SECURITY          │
                   EVENT            ▼
                              AGENT WALLET
                                   │
                              SOLANA SIGN
                                   │
                               SETTLEMENT
                                   │
                              VERIFICATION
                                   │
                                LEDGER
                              ┌────┴────┐
                            DEBIT     CREDIT
                              │          │
                              └────┬─────┘
                                   │
                              TRANSACTION
                                   │
                            ACTIVITY/AUDIT
