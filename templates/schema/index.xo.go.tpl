{{- $i := .Data -}}
// {{ func_name_context $i }} retrieves a row from '{{ schema $i.Table.SQLName }}' as a {{ $i.Table.GoName }}.
//
// Generated from index '{{ $i.SQLName }}'.
{{ func_context $i }} {
	// query
	{{ sqlstr "index" $i }}
	// run
	logf(sqlstr, {{ params $i.Fields false }})
{{- if $i.IsUnique }}
	{{ short $i.Table }} := {{ $i.Table.GoName }}{
	{{- if $i.Table.PrimaryKeys }}
		_exists: true,
	{{ end -}}
	}
	qb := sqrl.Expr(sqlstr, {{  params $i.Fields false }})
	if err := db.QueryRowContext(ctx, &{{short $i.Table}}, qb); err != nil {
		return nil, logerror(err)
	}
	return &{{ short $i.Table }}, nil
{{- else }}
	// process
	var res []*{{ $i.Table.GoName }}
	qb  := sqrl.Expr(sqlstr, {{  params $i.Fields false }})
	if err := db.QueryContext(ctx, &res, qb); err != nil {
		return nil, logerror(err)
	}
	return res, nil
{{- end }}
}

{{ if context_both -}}
// {{ func_name $i }} retrieves a row from '{{ schema $i.Table.SQLName }}' as a {{ $i.Table.GoName }}.
//
// Generated from index '{{ $i.SQLName }}'.
{{ func $i }} {
	return {{ func_name_context $i }}({{ names "" "context.Background()" "db" $i }})
}
{{- end }}

