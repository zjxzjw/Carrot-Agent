import { Table, TableProps } from 'antd'

interface DataTableProps<T extends object> extends TableProps<T> {
  dataSource: T[]
  rowKey?: string | ((record: T) => string)
  pagination?: false | {
    current?: number
    pageSize?: number
    total?: number
    onChange?: (page: number, pageSize: number) => void
  }
}

export function DataTable<T extends object>({ 
  dataSource, 
  rowKey = 'id',
  pagination = { pageSize: 10 },
  ...props 
}: DataTableProps<T>) {
  return (
    <Table<T>
      dataSource={dataSource}
      rowKey={rowKey}
      pagination={pagination}
      {...props}
    />
  )
}

export default DataTable