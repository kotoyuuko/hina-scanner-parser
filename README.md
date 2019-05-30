# 北科编译原理实验

## 简介

北京科技大学计算机科学与技术专业编译原理实验，要求实现词法分析和语法分析功能。

代码使用 Go 语言实现，解析一个名为 Hina 的简单脚本语言。

## 词法分析

## 语法分析

### EBNF of Hina

    program   ::= {statement}<EOF>;
    statement ::= define | assign | func | funccall | if;
    define    ::= <LET> <IDENT> <ASSIGN> expr;
    assign    ::= <IDENT> <ASSIGN> expr;
    expr      ::= func | funccall | <IDENT> | <INT> | <FLOAT> | <CHAR> | <STRING> | calc;
    operator  ::= <ADD> | <SUB> | <MUL> | <QUO> | <REM> | <AND> | <OR> | <XOR> | <SHL> | <SHR>
    condition ::= 
    calc      ::= <IDENT>operator[<IDENT>]
    func      ::= <LPAREN>[<IDENT>{<COMMA><IDENT>}]<RPAREN><LBRACE>program<RBRACE>;
    funcall   ::= <IDENT><LPAREN>[expr{,expr}]<RPAREN>;
    if        ::= 
